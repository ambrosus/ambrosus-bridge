import os
import re
import time
from collections import namedtuple
from dataclasses import dataclass
from pathlib import Path
from pprint import pprint

import requests

LOGS_FILE = Path.cwd() / "files" / "logs.txt"
LAST_LOG_FILE = Path.cwd() / "files" / "last_log.txt"

LOKI_URL = os.getenv("LOKI_URL")  # ex: 'loki:QWE123qwe123qwe123qwe@loki.ambrosus-ops.io'
LOKI_PARAMS = os.getenv("LOKI_PARAMS")  # ex: '{host="dev-parity0", container_name="parity"}'
LOKI_LIMIT = 5000

LOG_RE = re.compile(r'epoch \((\d+), .* #(\d+)')


@dataclass
class Log:
    timestamp: int
    finalize_at: int
    signaled_at: int


def main():
    last_log = read_last_log()
    print("Start working from", last_log)

    while True:
        try:
            logs = fetch_logs(last_log.timestamp)
        except Exception as e:
            print("Exception while fetching logs", e)
            time.sleep(60)
            continue
        print(f"fetched {len(logs)} logs")

        finalize_logs = filter(None, map(parse_loki_log, logs))  # filter out not finalize logs
        finalize_logs = list(filter(
            lambda log, last_log_=last_log: log.finalize_at < last_log_.finalize_at,
            finalize_logs))  # filter out already saved logs

        last_log.timestamp = int(logs[-1][0])  # set last_log timestamp to last log timestamp 🤯

        if finalize_logs:
            # set last_log last block to last finalize log block to avoid doubles
            last_log.finalize_at = finalize_logs[-1].finalize_at
            print(f"find {len(logs)} new logs, latest is {last_log}")
            save_new_logs(logs)

        if len(logs) < LOKI_LIMIT:
            print(f"Fetched less logs than limit ({len(logs)} / {LOKI_LIMIT}); Sleep 60 sec")
            time.sleep(60)


def fetch_logs(from_timestamp: int) -> [Log]:
    print(f"fetching logs {from_timestamp=}")
    r = requests.get(f"https://{LOKI_URL}/loki/api/v1/query_range", params={
        "query": LOKI_PARAMS,
        "start": from_timestamp,
        "limit": LOKI_LIMIT,
        "direction": "forward",
    })
    data_result = r.json()["data"]["result"]
    if not data_result:
        return []
    return data_result[0]["values"]


def parse_loki_log(log: [str, str]) -> Log:
    """
    :param log looks like this
    ['1657798570283477710',

    '2022-07-14 11:36:10 UTC Verifier #0 TRACE engine  extracting epoch validator set'
    'for epoch (216781, 0x222c…b6b0) signalled at #216776']
    """
    regex_match = LOG_RE.findall(log[1])
    if not regex_match:
        return None

    timestamp = log[0]
    finalize_at, signaled_at = regex_match[0]
    return Log(int(timestamp), int(finalize_at), int(signaled_at))


def save_new_logs(logs: [Log]):
    logs_str = '\n'.join([
        f"{log.finalize_at} {log.signaled_at}"
        for log in logs
    ])
    with LOGS_FILE.open("a") as f:
        f.write(logs_str)

    last_log_str = f"{logs[-1].timestamp} {logs[-1].finalize_at}"
    LAST_LOG_FILE.write_text(last_log_str)


def read_last_log() -> Log:
    try:
        log = LAST_LOG_FILE.read_text()
        timestamp, finalize_at = map(int, log.split(" "))
        return Log(timestamp, finalize_at, signaled_at=0)
    except FileNotFoundError:
        LAST_LOG_FILE.parent.mkdir(exist_ok=True)
        LAST_LOG_FILE.write_text(f"{time.time_ns()} 0")
        print(f"For first launch, set wanted timestamp (first number) in {LAST_LOG_FILE} file. "
              f"\n It's now set to current timestamp")
        exit(1)
    except Exception as e:
        raise Exception(f"Error when read {LAST_LOG_FILE} file") from e


main()
