import asyncio
import os

import yt_dlp

from injector import inject

from config import Config
from downloader import Downloader
from yt_dlp_logger import YtDlpLogger


class TwitchDownloader(Downloader):
    @inject
    def __init__(self, config: Config, yt_dlp_logger: YtDlpLogger):
        self._config = config
        self._yt_dlp_logger = yt_dlp_logger

    async def download(self, url: str) -> None:
        if not url:
            raise ValueError("url is required")

        def sync():
            ydl_opts = {
                "logger": self._yt_dlp_logger,
                "format": "best",
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"],
                    "[%(uploader)s] %(title)s.%(ext)s",
                ),
                # --- The Retry Trio ---
                "retries": 10,  # Generic network retries
                "fragment_retries": 10,  # Video chunk/fragment retries
                "extractor_retries": 10,  # Website parsing/scraping retries
                "file_access_retries": 10,  # Local disk/NAS access retries
                # Sleep between each retry attempt (all retry types)
                "retry_sleep_functions": {
                    "http": lambda _: 15,
                    "fragment": lambda _: 15,
                    "file_access": lambda _: 15,
                    "extractor": lambda _: 15,
                },
                # --- Precise Timing Control ---
                "sleep_interval": 15,  # Seconds to wait between download tasks
                "max_sleep_interval": 15,  # Keep it strictly at 15s (no randomization)
                "sleep_requests": 5,  # Wait 5s between finding info for each video
                # --- Safety Buffers ---
                "socket_timeout": 30,  # Wait 30s before considering a socket "dead"
                # Suppress FFmpeg's direct stderr output (it bypasses yt-dlp's logger)
                "postprocessor_args": {"ffmpeg": ["-loglevel", "warning"]},
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)
