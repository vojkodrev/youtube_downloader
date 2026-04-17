import asyncio
import os
import re

from loguru import logger

from injector import inject

from config import Config
from downloader_map import DownloaderMap

SERVICE_TO_DOWNLOADER = {
    "youtube": "youtube_video",
    "twitch": "twitch",
    "rumble": "rumble",
}

DOWNLOAD_FILE_RE = re.compile(r"^video\.([^.]+)\.([^.]+)\.download$")


class DownloadVideoRequestPoller:
    @inject
    def __init__(self, config: Config, downloaders: DownloaderMap):
        self._config = config
        self._downloaders = downloaders

    async def poll(self) -> None:
        while True:
            folder = self._config["output_folder"]
            try:
                entries = os.listdir(folder)
            except Exception as e:
                logger.bind(streamer="-").error(f"Failed to list streams folder: {e}")
                await asyncio.sleep(60)
                continue

            for name in entries:
                m = DOWNLOAD_FILE_RE.match(name)
                if not m:
                    continue

                service, video_id = m.group(1), m.group(2)
                path = os.path.join(folder, name)
                log = logger.bind(streamer=f"{service}/{video_id}")

                try:
                    downloader_key = SERVICE_TO_DOWNLOADER.get(service)
                    if downloader_key is None:
                        log.warning(f"Unknown service '{service}', skipping {name}")
                    else:
                        url = open(path).read().strip()
                        log.info(f"Downloading {url}")
                        await self._downloaders[downloader_key].download(url, f"{service}/{video_id}")
                        log.info("Download finished")
                except Exception as e:
                    log.error(f"Failed to process video: {e}")
                finally:
                    os.remove(path)
                    log.info("Request file removed")

            await asyncio.sleep(60)
