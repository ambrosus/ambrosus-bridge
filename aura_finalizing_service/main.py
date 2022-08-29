import os
import re
import time
from dataclasses import dataclass
from pathlib import Path

import requests

LOGS_FILE = Path.cwd() / "files" / "logs.bin"
LAST_LOG_TIMESTAMP_FILE = Path.cwd() / "files" / "last_log_timestamp.txt"

LOKI_URL = os.getenv("LOKI_URL")  # ex: 'loki:QWE123qwe123qwe123qwe@loki.ambrosus-ops.io'
LOKI_PARAMS = os.getenv("LOKI_PARAMS")  # ex: '{host="dev-parity0", container_name="parity"}'
LOKI_LIMIT = 5000

LOG_RE = re.compile(r'epoch \((\d+), .* #(\d+)')


@dataclass
class Log:
    timestamp: int
    finalize_at: int
    signaled_at: int

    def to_binary(self):
        # block numbers is uint64 => 8 bytes
        return self.signaled_at.to_bytes(8, byteorder='little') + \
               self.finalize_at.to_bytes(8, byteorder='little')


def main():
    last_log_timestamp = read_last_log()
    print("Start working from", last_log_timestamp)

    while True:
        try:
            logs = fetch_logs(last_log_timestamp)
        except Exception as e:
            print("Exception while fetching logs", e)
            time.sleep(60)
            continue
        if logs:
            last_log_timestamp = int(logs[-1][0])
        print(f"fetched {len(logs)} logs, latest is {last_log_timestamp}")

        finalize_logs = list(filter(None, map(parse_loki_log, logs)))  # filter out not finalize logs
        if finalize_logs:
            print(f"find {len(finalize_logs)} new finalize logs")
            save_new_logs(finalize_logs)

        save_last_log(last_log_timestamp)

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
    }, timeout=5)
    data_result = r.json()["data"]["result"]
    if not data_result:
        return []
    return [j for i in data_result for j in i["values"]]


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
    with LOGS_FILE.open("ab") as f:
        for log in logs:
            f.write(log.to_binary())


def save_last_log(last_log_timestamp: int):
    LAST_LOG_TIMESTAMP_FILE.write_text(str(last_log_timestamp))


def read_last_log() -> int:
    try:
        return int(LAST_LOG_TIMESTAMP_FILE.read_text())
    except FileNotFoundError:
        LAST_LOG_TIMESTAMP_FILE.parent.mkdir(exist_ok=True)
        LAST_LOG_TIMESTAMP_FILE.write_text(str(time.time_ns()))
        print(f"For first launch, set wanted timestamp (in ns) in {LAST_LOG_TIMESTAMP_FILE} file. "
              f"\n It's now set to current timestamp")
        exit(1)
    except Exception as e:
        raise Exception(f"Error when read {LAST_LOG_TIMESTAMP_FILE} file") from e


if __name__ == "__main__":
    main()
