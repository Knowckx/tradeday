import json
from pathlib import Path
from datetime import date, timedelta
import sys

import pandas_market_calendars as pmc


DEFAULT_TRUTH_TABLE_PATH = Path(__file__).resolve().parent.parent / "testdata" / "us_stock_truth_table.json"
US_STOCK_CALENDAR_NAME = "NYSE"


def load_truth_table(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def build_trade_day_set(calendar_name: str, start: str, end: str) -> set[str]:
    calendar = pmc.get_calendar(calendar_name)
    schedule = calendar.schedule(start_date=start, end_date=end)
    return {idx.strftime("%Y-%m-%d") for idx in schedule.index}


def compare_truth_table(payload: dict, calendar_name: str) -> list[tuple[str, bool, bool]]:
    trade_days = build_trade_day_set(calendar_name, payload["start"], payload["end"])
    mismatches: list[tuple[str, bool, bool]] = []

    current = date.fromisoformat(payload["start"])
    end = date.fromisoformat(payload["end"])
    truth_days = payload["days"]

    while current <= end:
        day = current.isoformat()
        got = day in trade_days
        want = truth_days[day]
        if got != want:
            mismatches.append((day, got, want))

        current += timedelta(days=1)

    return mismatches


def main() -> int:
    if len(sys.argv) > 2:
        raise SystemExit("usage: python verify_us_stock_truth_table_with_pandas_market_calendars.py [truth_table_path]")

    truth_table_path = DEFAULT_TRUTH_TABLE_PATH
    if len(sys.argv) == 2:
        truth_table_path = Path(sys.argv[1]).resolve()

    payload = load_truth_table(truth_table_path)
    mismatches = compare_truth_table(payload, US_STOCK_CALENDAR_NAME)

    print(f"truth_table={truth_table_path}")
    print(f"calendar_id={payload['calendar_id']}")
    print(f"range={payload['start']}..{payload['end']}")
    print(f"mismatch_count={len(mismatches)}")

    for day, got, want in mismatches[:50]:
        print(f"{day} pandas_market_calendars={got} truth_table={want}")

    return 0 if not mismatches else 1


if __name__ == "__main__":
    raise SystemExit(main())
