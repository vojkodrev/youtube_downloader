import asyncio
import itertools
import os
import sys
import tomllib
from dotenv import load_dotenv
from googleapiclient.discovery import build
from loguru import logger
import yt_dlp


async def get_live_video_id(api_keys, channel_title=None, channel_id=None):
    if not channel_title and not channel_id:
        raise ValueError("at least one of channel_title or channel_id is required")

    def get_live_video_id_sync():
        youtube = build("youtube", "v3", developerKey=api_keys.next())

        if channel_id:
            request = youtube.search().list(
                part="snippet",
                channelId=channel_id,
                eventType="live",
                type="video",
                maxResults=1,
            )
            response = request.execute()
            if response.get("items"):
                return response["items"][0]["id"].get("videoId")
        else:
            request = youtube.search().list(
                part="snippet",
                q=channel_title,
                eventType="live",
                type="video",
                maxResults=1,
            )
            response = request.execute()
            if response.get("items"):
                channel_title_resp = response["items"][0]["snippet"]["channelTitle"]
                if channel_title_resp.lower() == channel_title.lower():
                    return response["items"][0]["id"].get("videoId")

        return None

    loop = asyncio.get_running_loop()
    return await loop.run_in_executor(None, get_live_video_id_sync)


async def get_channel_title(api_keys, channel_id):
    if not channel_id:
        raise ValueError("channel_id is required")

    def get_channel_title_sync():
        youtube = build("youtube", "v3", developerKey=api_keys.next())
        response = youtube.channels().list(part="snippet", id=channel_id).execute()
        if not response.get("items"):
            raise ValueError(f"Could not find channel with ID '{channel_id}'")
        return response["items"][0]["snippet"]["title"]

    loop = asyncio.get_running_loop()
    return await loop.run_in_executor(None, get_channel_title_sync)


def get_video_url(video_id):
    if not video_id:
        raise ValueError("video_id is required")
    return f"https://www.youtube.com/watch?v={video_id}"


async def download_live_from_start(url, download_folder="."):
    if not url:
        raise ValueError("url is required")

    def download_live_from_start_sync():
        ydl_opts = {
            "format": "bestvideo+bestaudio/best",
            # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
            "live_from_start": True,
            "merge_output_format": "mp4",
            "outtmpl": os.path.join(download_folder, "%(title)s.%(ext)s"),
            # "http_headers": {"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"},
            # Optional: retry if the stream connection drops
            # "ignoreerrors": True,
            # "concurrent_fragment_downloads": 10,  # Download 10 chunks at once
        }

        with yt_dlp.YoutubeDL(ydl_opts) as ydl:
            ydl.download([url])

    loop = asyncio.get_running_loop()
    await loop.run_in_executor(None, download_live_from_start_sync)


class FibonacciSleep:
    SHORT = [5, 8, 13, 21, 30]
    LONG = [20, 30, 45]

    def __init__(self, intervals=None):
        self._intervals = intervals or self.SHORT
        self._index = 0

    async def sleep(self):
        interval = self._intervals[self._index]
        self._index = min(self._index + 1, len(self._intervals) - 1)
        await asyncio.sleep(interval * 60)
        return interval

    def peek(self):
        return self._intervals[self._index]

    def reset(self):
        self._index = 0


async def poll_and_download(api_keys, channel_title=None, channel_id=None, download_folder="."):
    if channel_id and not channel_title:
        channel_title = await get_channel_title(api_keys, channel_id)

    identifier = channel_title or channel_id
    log = logger.bind(streamer=identifier)

    log.info(f"Resolved channel ID '{channel_id}'. Polling started...")

    sleep_offline = FibonacciSleep(FibonacciSleep.LONG)
    sleep_err = FibonacciSleep(FibonacciSleep.SHORT)

    while True:
        try:
            video_id = await get_live_video_id(api_keys, channel_title, channel_id)

            if video_id:
                sleep_offline.reset()
                url = get_video_url(video_id)
                log.info(f"Streamer is LIVE! Downloading from: {url}")
                await download_live_from_start(url, download_folder)
                sleep_err.reset()
                log.info("Download finished. Resuming poll...")
            else:
                log.info(f"Streamer is offline. Checking again in {sleep_offline.peek()} minutes...")
                await sleep_offline.sleep()
        except Exception as e:
            log.error(f"Error: {e}. Retrying in {sleep_err.peek()} minutes...")
            await sleep_err.sleep()


class ApiKeyPool:
    def __init__(self):
        keys = [k.strip() for k in os.getenv("API_KEYS", "").split(",") if k.strip()]
        if not keys:
            raise RuntimeError("No API keys found in API_KEYS env var")
        self._cycle = itertools.cycle(keys)

    def next(self):
        return next(self._cycle)


def main():
    load_dotenv(os.path.join(os.path.dirname(__file__), ".env"))
    api_keys = ApiKeyPool()

    with open(os.path.join(os.path.dirname(__file__), "config.toml"), "rb") as f:
        config = tomllib.load(f)

    channel_ids = config["channel_ids"]
    output_folder = config["output_folder"]
    log_format = config["log_format"]

    logger.remove()
    logger.add(sys.stderr, format=log_format)

    if not os.path.isdir(output_folder):
        logger.bind(streamer="-").error(
            f"Output folder does not exist: {output_folder}"
        )
        exit(1)

    async def poll_all_channels():
        await asyncio.gather(
            *[
                poll_and_download(api_keys, channel_id=cid, download_folder=output_folder)
                for cid in channel_ids
            ]
        )

    asyncio.run(poll_all_channels())


if __name__ == "__main__":
    main()
