from dataclasses import dataclass


@dataclass
class VideoMetadata:
    title: str | None
    live_start_time: str | None
