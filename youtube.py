import argparse
import os
import time
from dotenv import load_dotenv
from googleapiclient.discovery import build
from loguru import logger
import yt_dlp


def get_live_video_id(channel_title=None, channel_id=None):
    if not channel_title and not channel_id:
        raise ValueError("at least one of channel_title or channel_id is required")

    youtube = build("youtube", "v3", developerKey=os.getenv("API_KEY"))

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


def get_channel_title(channel_id):
    youtube = build("youtube", "v3", developerKey=os.getenv("API_KEY"))
    response = youtube.channels().list(part="snippet", id=channel_id).execute()
    if not response.get("items"):
        raise ValueError(f"Could not find channel with ID '{channel_id}'")
    return response["items"][0]["snippet"]["title"]


def get_video_url(video_id):
    return f"https://www.youtube.com/watch?v={video_id}"


def download_live_from_start(url, download_folder="."):
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


def poll_and_download(channel_title=None, channel_id=None, download_folder="."):
    FIBONACCI_INTERVALS = [5, 8, 13, 21, 30]

    if channel_id and not channel_title:
        channel_title = get_channel_title(channel_id)
        logger.info(f"Resolved channel ID '{channel_id}' to '{channel_title}'")

    identifier = channel_title or channel_id

    logger.info(
        f"Polling for '{identifier}' with Fibonacci backoff {FIBONACCI_INTERVALS} minutes..."
    )

    fib_index = 0

    while True:
        try:
            video_id = get_live_video_id(channel_title, channel_id)

            if video_id:
                fib_index = 0
                url = get_video_url(video_id)
                logger.info(f"Streamer is LIVE! Downloading from: {url}")
                download_live_from_start(url, download_folder)
                logger.info("Download finished. Resuming poll...")
            else:
                interval = FIBONACCI_INTERVALS[fib_index]
                logger.info(
                    f"Streamer is offline. Checking again in {interval} minutes..."
                )
                time.sleep(interval * 60)
                fib_index = min(fib_index + 1, len(FIBONACCI_INTERVALS) - 1)
        except Exception as e:
            logger.error(f"Error: {e}. Retrying in 1 minute...")
            time.sleep(60)


def main():
    load_dotenv()
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "-ct",
        "--channel-title",
        dest="channel_title",
        help="YouTube channel title to watch",
    )
    parser.add_argument(
        "-ci", "--channel-id", dest="channel_id", help="YouTube channel ID to watch"
    )
    parser.add_argument(
        "-o",
        "--output",
        default=".",
        help="Folder to save downloads (default: current directory)",
    )
    args = parser.parse_args()

    if not args.channel_title and not args.channel_id:
        parser.error("at least one of --channel-title or --channel-id is required")

    if not os.path.isdir(args.output):
        logger.error(f"Output folder does not exist: {args.output}")
        exit(1)

    poll_and_download(args.channel_title, args.channel_id, args.output)


if __name__ == "__main__":
    main()
