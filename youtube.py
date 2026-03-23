import os
from dotenv import load_dotenv
from googleapiclient.discovery import build

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
        # if channel_title.lower() == streamer_name.lower():
        return response["items"][0]["id"]["videoId"]

    return None


def get_video_url(video_id):
    return f"https://www.youtube.com/watch?v={video_id}"


def download_live_from_start(url):
    ydl_opts = {
        'format': 'bestvideo+bestaudio/best',
        # CRITICAL: This flag tells yt-dlp to start from the beginning of the DVR
        'live_from_start': True,
        'merge_output_format': 'mp4',
        'outtmpl': '%(title)s.%(ext)s',
        # Optional: retry if the stream connection drops
        'ignoreerrors': True,
    }

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        ydl.download([url])


if __name__ == "__main__":
    video_id = get_live_video_id(STREAMER_NAME)
    if video_id:
        url = get_video_url(video_id)
        print(f"Streamer is LIVE! Watch here: {url}")
        download_live_from_start(url)
    else:
        print("Streamer is offline.")
