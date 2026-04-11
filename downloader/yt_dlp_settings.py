def _retry_sleep(n):  # noqa: ARG001
    return 15


SHARED_YT_DLP_SETTINGS = {
    # --- The Retry Trio ---
    "retries": 10,  # Generic network retries
    "fragment_retries": 10,  # Video chunk/fragment retries
    "extractor_retries": 10,  # Website parsing/scraping retries
    "file_access_retries": 10,  # Local disk/NAS access retries
    # Sleep between each retry attempt (all retry types)
    "retry_sleep_functions": {
        "http": _retry_sleep,
        "fragment": _retry_sleep,
        "file_access": _retry_sleep,
        "extractor": _retry_sleep,
    },
    # --- Precise Timing Control ---
    "sleep_interval": 15,  # Seconds to wait between download tasks
    "max_sleep_interval": 15,  # Keep it strictly at 15s (no randomization)
    "sleep_requests": 5,  # Wait 5s between finding info for each video
    # --- Safety Buffers ---
    "socket_timeout": 30,  # Wait 30s before considering a socket "dead"
}
