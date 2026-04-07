import itertools
import os


class YoutubeApiKeyPool:
    def __init__(self):
        keys = [k.strip() for k in os.getenv("API_KEYS", "").split(",") if k.strip()]
        if not keys:
            raise RuntimeError("No API keys found in API_KEYS env var")
        self._cycle = itertools.cycle(keys)

    def next(self):
        return next(self._cycle)
