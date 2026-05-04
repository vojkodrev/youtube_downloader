from __future__ import annotations

import re
import os
from datetime import datetime, timezone

import ffmpeg
from injector import inject
from loguru import logger

from config import Config
from metadata_provider_map import MetadataProviderMap
from video_validator import VideoValidator


class DownloadedVideoValidator(VideoValidator):
    @inject
    def __init__(self, meta_providers: MetadataProviderMap, config: Config):
        self._meta_providers = meta_providers
        self._config = config

    async def validate(self, video_id: str, mode: str, streamer: str = "") -> list[str]:
        log = logger.bind(streamer=streamer)
        log.debug(f"Validating video '{video_id}' (mode={mode})")

        metadata = await self._meta_providers[mode].get_video_metadata(video_id)
        if not metadata:
            log.debug(f"No metadata found for video '{video_id}', skipping validation")
            return []
        if not metadata.title:
            log.debug(f"Metadata has no title for video '{video_id}', skipping validation")
            return []

        log.debug(f"Metadata fetched: title='{metadata.title}', duration_seconds={metadata.duration_seconds}, live_start_time={metadata.live_start_time}")

        if metadata.duration_seconds:
            stream_duration = metadata.duration_seconds
            log.debug(f"Using contentDetails duration: {stream_duration:.0f}s")
        elif metadata.live_start_time:
            start_time = datetime.fromisoformat(metadata.live_start_time.replace("Z", "+00:00"))
            stream_duration = (datetime.now(timezone.utc) - start_time).total_seconds()
            log.debug(f"Using live_start_time '{metadata.live_start_time}' to compute stream duration: {stream_duration:.0f}s")
        else:
            log.debug(f"Cannot determine stream duration for video '{video_id}', skipping validation")
            return []

        output_folder = self._config["output_folder"]
        log.debug(f"Looking for recorded files in '{output_folder}' matching title '{metadata.title}'")
        recorded_duration = self._get_recorded_duration(metadata.title, output_folder, log)
        if recorded_duration is None:
            log.debug(f"No recorded file found for title '{metadata.title}', skipping validation")
            return []

        gap = stream_duration - recorded_duration
        log.debug(f"Stream duration: {stream_duration:.0f}s, recorded duration: {recorded_duration:.0f}s, gap: {gap:.0f}s")

        if gap <= 300:
            msg = f"Recorded duration ({recorded_duration:.0f}s) is within 5 minutes of stream duration ({stream_duration:.0f}s). Download not needed."
            log.info(msg)
            return [msg]

        log.debug(f"Gap of {gap:.0f}s exceeds threshold, download needed")
        return []

    def _title_to_pattern(self, title: str) -> str:
        return r".*".join(re.escape(w) for w in title.split(" "))

    def _get_recorded_duration(self, title: str, output_folder: str, log) -> float | None:
        try:
            files = os.listdir(output_folder)
        except OSError as e:
            log.warning(f"Could not list output folder '{output_folder}': {e}")
            return None

        log.debug(f"Found {len(files)} file(s) in output folder")
        title_pat = self._title_to_pattern(title)

        # Source mp4 takes priority (partXX files exist only while splitting is in progress)
        base_pattern = re.compile(title_pat + r"\.\w+$", re.IGNORECASE)
        base_files = [f for f in files if base_pattern.search(f) and ".temp." not in f]
        if base_files:
            log.debug(f"Found source file(s): {base_files}, using '{base_files[0]}'")
            return self._get_duration(os.path.join(output_folder, base_files[0]), log)

        # No source file — splitting is done, sum the part files
        part_pattern = re.compile(title_pat + r" part\d{2}\.\w+$", re.IGNORECASE)
        parts = sorted(f for f in files if part_pattern.search(f) and ".temp." not in f)
        if not parts:
            log.debug(f"No source or part files matched title '{title}'")
            return None

        log.debug(f"Found {len(parts)} part file(s): {parts}")
        return self._sum_durations(parts, output_folder, log)

    def _sum_durations(self, filenames: list[str], folder: str, log) -> float | None:
        total = 0.0
        for filename in filenames:
            duration = self._get_duration(os.path.join(folder, filename), log)
            if duration is None:
                log.warning(f"Could not read duration for '{filename}', aborting sum")
                return None
            log.debug(f"  '{filename}': {duration:.0f}s")
            total += duration
        log.debug(f"Total duration from {len(filenames)} part(s): {total:.0f}s")
        return total

    def _get_duration(self, path: str, log) -> float | None:
        try:
            probe = ffmpeg.probe(path)
            duration = float(probe["format"]["duration"])
            log.debug(f"Probed '{os.path.basename(path)}': {duration:.0f}s")
            return duration
        except (ffmpeg.Error, KeyError, ValueError) as e:
            log.warning(f"ffmpeg probe failed for '{path}': {e}")
            return None
