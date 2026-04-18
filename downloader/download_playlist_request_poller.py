import asyncio
import os
import re
from urllib.parse import urlparse, parse_qs

from loguru import logger

from injector import inject

from config import Config
from downloader_map import DownloaderMap
from metadata_provider_map import MetadataProviderMap

SERVICE_TO_DOWNLOADER = {
    "youtube": "youtube_video",
}

DOWNLOAD_FILE_RE = re.compile(r"^playlist\.([^.]+)\.([^.]+)\.download$")


class DownloadPlaylistRequestPoller:
    @inject
    def __init__(
        self,
        config: Config,
        downloaders: DownloaderMap,
        metadata_providers: MetadataProviderMap,
    ):
        self._config = config
        self._downloaders = downloaders
        self._metadata_providers = metadata_providers

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

                service = m.group(1)
                path = os.path.join(folder, name)
                log = logger.bind(streamer=service)

                try:
                    downloader_key = SERVICE_TO_DOWNLOADER.get(service)
                    if downloader_key is None:
                        log.warning(f"Unknown service '{service}', skipping {name}")
                        continue

                    url = open(path).read().strip()
                    playlist_id = parse_qs(urlparse(url).query).get("list", [None])[0]
                    if not playlist_id:
                        log.error(f"Could not extract playlist id from {url}")
                    else:
                        log.info(f"Fetching playlist {url}")
                        metadata_provider = self._metadata_providers[downloader_key]
                        videos = await metadata_provider.get_playlist_videos(
                            playlist_id
                        )
                        for video in videos:
                            video_url = metadata_provider.get_video_url(
                                video["video_id"]
                            )
                            video_log = log.bind(
                                streamer=f"{service}/{video['video_id']}"
                            )
                            video_log.info(f"Downloading {video_url}")
                            try:
                                await self._downloaders[downloader_key].download(
                                    video_url, f"{service}/{video['video_id']}"
                                )
                                video_log.info("Download finished")
                            except Exception as e:
                                video_log.error(f"Download failed: {e}")
                except Exception as e:
                    log.error(f"Failed to process playlist: {e}")
                finally:
                    os.remove(path)
                    log.info("Request file removed")

            await asyncio.sleep(60)
