from fibonacci_sleep import FibonacciSleep


class FibonacciSleepFactory:
    def create(self, interval_type: str, mode: str) -> FibonacciSleep:
        if mode == "twitch":
            return FibonacciSleep(intervals=[15])
        if mode == "youtube_live":
            if interval_type == "short":
                return FibonacciSleep(intervals=[5, 8, 13, 21, 34])
            if interval_type == "long":
                return FibonacciSleep(intervals=[13, 21, 34, 55])
            raise ValueError(f"Unsupported interval_type: {interval_type}")
        raise ValueError(f"Unsupported mode: {mode}")
