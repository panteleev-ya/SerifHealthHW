import json
import time
import logging
import os
import ijson


log_file = f"{os.path.basename(__file__).split('.')[0]}.log"
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
    handlers=[
        logging.FileHandler(log_file),
        logging.StreamHandler()
    ]
)
logger = logging.getLogger(__name__)


def any_match(s, s_filter):
    for option in s_filter:
        if option in s:
            return True
    return False


def all_match(s, filters):
    for s_filter in filters:
        if not any_match(s, s_filter):
            return False
    return True


def extract_filtered_links(file_path, description_filters=None):
    if description_filters is None:
        description_filters = []
    filtered_links = []
    with open(file_path, 'r') as file:
        network_files = ijson.items(file, 'reporting_structure.item.in_network_files.item')
        for network_file in network_files:
            description_lc = str(network_file.get("description", "")).lower()
            if all_match(description_lc, description_filters):
                filtered_links.append(str(network_file["location"]))
    return filtered_links


if __name__ == '__main__':
    file_path = '../2024-08-01_anthem_index.json'
    filters = [["new york", "ny"], ["ppo"]]

    start_time = time.perf_counter()

    logger.info("Started extracting...")
    links = extract_filtered_links(file_path, description_filters=filters)
    json.dump(links, open('output.json', 'w'))

    end_time = time.perf_counter()
    elapsed_time = end_time - start_time

    logger.info(f"Filtered links found: {len(links)}")
    logger.info(f"Extracting completed in {elapsed_time} seconds")

