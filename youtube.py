import argparse
import os
import time
from dotenv import load_dotenv
from googleapiclient.discovery import build
from loguru import logger
import yt_dlp

def get_live_video_id(streamer_name):
    youtube = build("youtube", "v3", developerKey=os.getenv("API_KEY"))

    request = youtube.search().list(
        part="snippet", q=streamer_name, eventType="live", type="video", maxResults=1
    )
    response = request.execute()

    if response.get("items"):
        channel_title = response["items"][0]["snippet"]["channelTitle"]
        if channel_title.lower() == streamer_name.lower():
            return response["items"][0]["id"].get("videoId")

    return None


def get_video_url(video_id):
    return f"https://www.youtube.com/watch?v={video_id}"


def download_live_from_start(url):
    ydl_opts = {
        "format": "bestvideo+bestaudio/best",
        # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
        "live_from_start": True,
        "merge_output_format": "mp4",
        "outtmpl": "%(title)s.%(ext)s",
        # Optional: retry if the stream connection drops
        "ignoreerrors": True,
    }

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])


def poll_and_download(streamer_name, interval_minutes=15):
    logger.info(f"Polling every {interval_minutes} minutes for '{streamer_name}'...")
    while True:
        video_id = get_live_video_id(streamer_name)
        if video_id:
            url = get_video_url(video_id)
            logger.info(f"Streamer is LIVE! Downloading from: {url}")
            download_live_from_start(url)
            logger.info("Download finished. Resuming poll...")
        else:
            logger.info(f"Streamer is offline. Checking again in {interval_minutes} minutes...")
        time.sleep(interval_minutes * 60)


if __name__ == "__main__":
    load_dotenv()
    parser = argparse.ArgumentParser()
    parser.add_argument("streamer_name", help="YouTube channel name to watch")
    args = parser.parse_args()
    poll_and_download(args.streamer_name)
