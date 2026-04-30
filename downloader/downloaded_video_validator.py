from __future__ import annotations

import re
import os
from datetime import datetime, timezone

import ffmpeg
from injector import inject

from config import Config
from metadata_provider import MetadataProvider
from video_validator import VideoValidator


class DownloadedVideoValidator(VideoValidator):
    @inject
    def __init__(self, metadata_provider: MetadataProvider, config: Config):
        self._metadata_provider = metadata_provider
        self._config = config

    async def validate(self, video_id: str) -> bool:
        metadata = await self._metadata_provider.get_video_metadata(video_id)
        if not metadata or not metadata.title or not metadata.live_start_time:
            return True

        start_time = datetime.fromisoformat(metadata.live_start_time.replace("Z", "+00:00"))
        stream_duration = (datetime.now(timezone.utc) - start_time).total_seconds()

        output_folder = self._config["output_folder"]
        recorded_duration = self._get_recorded_duration(metadata.title, output_folder)
        if recorded_duration is None:
            return True

        return abs(stream_duration - recorded_duration) > 300

    def _get_recorded_duration(self, title: str, output_folder: str) -> float | None:
        try:
            files = os.listdir(output_folder)
        except OSError:
            return None

        part_pattern = re.compile(re.escape(title) + r" part\d+\.", re.IGNORECASE)
        parts = [f for f in files if part_pattern.search(f)]

        if parts:
            return self._sum_durations(parts, output_folder)

        base_pattern = re.compile(re.escape(title) + r"\.", re.IGNORECASE)
        base_files = [f for f in files if base_pattern.search(f)]
        if not base_files:
            return None

        return self._get_duration(os.path.join(output_folder, base_files[0]))

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
