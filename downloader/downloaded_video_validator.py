from __future__ import annotations

import re
import os
from datetime import datetime, timezone

import ffmpeg
from injector import inject

from config import Config
from metadata_provider_map import MetadataProviderMap
from video_validator import VideoValidator


class DownloadedVideoValidator(VideoValidator):
    @inject
    def __init__(self, meta_providers: MetadataProviderMap, config: Config):
        self._meta_providers = meta_providers
        self._config = config

    async def validate(self, video_id: str, mode: str) -> list[str]:
        metadata = await self._meta_providers[mode].get_video_metadata(video_id)
        if not metadata or not metadata.title:
            return []

        if metadata.duration_seconds:
            stream_duration = metadata.duration_seconds
        elif metadata.live_start_time:
            start_time = datetime.fromisoformat(metadata.live_start_time.replace("Z", "+00:00"))
            stream_duration = (datetime.now(timezone.utc) - start_time).total_seconds()
        else:
            return []

        output_folder = self._config["output_folder"]
        recorded_duration = self._get_recorded_duration(metadata.title, output_folder)
        if recorded_duration is None:
            return []

        gap = stream_duration - recorded_duration
        if gap <= 300:
            return [
                f"Recorded duration ({recorded_duration:.0f}s) is within 5 minutes of stream duration ({stream_duration:.0f}s). Download not needed."
            ]
        return []

    def _title_to_pattern(self, title: str) -> str:
        return r".*".join(re.escape(w) for w in title.split(" "))

    def _get_recorded_duration(self, title: str, output_folder: str) -> float | None:
        try:
            files = os.listdir(output_folder)
        except OSError:
            return None

        title_pat = self._title_to_pattern(title)

        # Source mp4 takes priority (partXX files exist only while splitting is in progress)
        base_pattern = re.compile(title_pat + r"\.\w+$", re.IGNORECASE)
        base_files = [f for f in files if base_pattern.search(f) and ".temp." not in f]
        if base_files:
            return self._get_duration(os.path.join(output_folder, base_files[0]))

        # No source file — splitting is done, sum the part files
        part_pattern = re.compile(title_pat + r" part\d{2}\.\w+$", re.IGNORECASE)
        parts = [f for f in files if part_pattern.search(f) and ".temp." not in f]
        if not parts:
            return None

        return self._sum_durations(parts, output_folder)

    def _sum_durations(self, filenames: list[str], folder: str) -> float | None:
        total = 0.0
        for filename in filenames:
            duration = self._get_duration(os.path.join(folder, filename))
            if duration is None:
                return None
            total += duration
        return total

    def _get_duration(self, path: str) -> float | None:
        try:
            probe = ffmpeg.probe(path)
            return float(probe["format"]["duration"])
        except (ffmpeg.Error, KeyError, ValueError):
            return None
