import itertools


class YoutubeApiKeyPool:
    def __init__(self, api_keys_str: str):
        keys = [k.strip() for k in api_keys_str.split(",") if k.strip()]
        if not keys:
            raise RuntimeError("No API keys found in API_KEYS env var")
        self._cycle = itertools.cycle(keys)

    def next(self):
        return next(self._cycle)
