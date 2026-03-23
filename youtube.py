import os
from dotenv import load_dotenv
from googleapiclient.discovery import build
import yt_dlp

load_dotenv()
API_KEY = os.getenv("API_KEY")
STREAMER_NAME = "The PrimeTime"


def get_live_video_id(streamer_name):
    youtube = build("youtube", "v3", developerKey=API_KEY)

    # Search for active live broadcasts for this channel
    request = youtube.search().list(
        part="snippet", q=streamer_name, eventType="live", type="video", maxResults=1
    )
    response = request.execute()

    if response.get("items"):
        channel_title = response["items"][0]["snippet"]["channelTitle"]
        if channel_title.lower() == streamer_name.lower():
            return response["items"][0]["id"]["videoId"]

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
    import time
    print(f"Polling every {interval_minutes} minutes for '{streamer_name}'...")
    while True:
        video_id = get_live_video_id(streamer_name)
        if video_id:
            url = get_video_url(video_id)
            print(f"Streamer is LIVE! Downloading from: {url}")
            download_live_from_start(url)
            print("Download finished. Resuming poll...")
        else:
            print("Streamer is offline. Checking again in 15 minutes...")
        time.sleep(interval_minutes * 60)


if __name__ == "__main__":
    poll_and_download(STREAMER_NAME)
