import asyncio


class FibonacciSleep:
    def __init__(self, intervals):
        self._intervals = intervals
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
