from video_validator import VideoValidator


class DummyVideoValidator(VideoValidator):
    async def validate(self, video_id: str) -> list[str]:
        return []
