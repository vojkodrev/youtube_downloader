import asyncio
import os
import sys

from loguru import logger
from injector import Injector

from ioc import Ioc
from multi_channel_poller import MultiChannelPoller


def main():
    ioc = Ioc()
    container = Injector([ioc])
    config = ioc._config

    logger.remove()
    logger.add(sys.stderr, format=config["log_format"])

    if not os.path.isdir(config["output_folder"]):
        logger.bind(streamer="-").error(
            f"Output folder does not exist: {config['output_folder']}"
        )
        exit(1)

    multi_poller = container.get(MultiChannelPoller)
    asyncio.run(multi_poller.poll_all())


if __name__ == "__main__":
    main()
