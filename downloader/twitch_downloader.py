import asyncio
import os

import yt_dlp

from injector import inject

from config import Config
from downloader import Downloader
from yt_dlp_logger import YtDlpLogger
from yt_dlp_settings import SHARED_YT_DLP_SETTINGS


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
                **SHARED_YT_DLP_SETTINGS,
                "logger": self._yt_dlp_logger,
                "format": "best",
                "merge_output_format": "mp4",
                "overwrites": True,
                "outtmpl": os.path.join(
                    self._config["output_folder"],
                    "[%(uploader)s] %(title)s.%(ext)s",
                ),
                # Suppress FFmpeg's direct stderr output (it bypasses yt-dlp's logger)
                "postprocessor_args": {"ffmpeg": ["-loglevel", "warning", "-stats", "-stats_period", "60"]},
                "external_downloader_args": {"ffmpeg": ["-loglevel", "warning", "-stats", "-stats_period", "60"]},
            }
            with yt_dlp.YoutubeDL(ydl_opts) as ydl:
                ydl.download([url])

        loop = asyncio.get_running_loop()
        await loop.run_in_executor(None, sync)
