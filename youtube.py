import os
from dotenv import load_dotenv
from googleapiclient.discovery import build

load_dotenv()
API_KEY = os.getenv("API_KEY")
STREAMER_NAME = "The PrimeTime"


def is_live(streamerName):
    youtube = build("youtube", "v3", developerKey=API_KEY)

    # Search for active live broadcasts for this channel
    request = youtube.search().list(
        part="snippet", q=streamerName, eventType="live", type="video", maxResults=1
    )
    response = request.execute()

    if response.get("items"):
        channel_title = response["items"][0]["snippet"]["channelTitle"]
        if True: #channel_title.lower() == streamerName.lower():
            video_id = response["items"][0]["id"]["videoId"]
            print(
                f"Streamer is LIVE! Watch here: https://www.youtube.com/watch?v={video_id}"
            )
            return True

    print("Streamer is offline.")
    return False


if __name__ == "__main__":
    is_live(STREAMER_NAME)
