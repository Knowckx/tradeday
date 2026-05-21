import json
from dataclasses import dataclass
from datetime import date
from pathlib import Path


ROOT = Path(__file__).resolve().parent.parent.parent


@dataclass(frozen=True)
class BitmapSpec:
    calendar_id: str
    truth_table_path: Path
    output_path: Path
    go_var_name: str
    comment: str


BITMAP_SPECS = {
    "cn_stock": BitmapSpec(
        calendar_id="cn_stock",
        truth_table_path=ROOT / "tests" / "testdata" / "cn_stock_truth_table.json",
        output_path=ROOT / "core" / "data" / "cn_days_bitmaps.go",
        go_var_name="CNStockTradeBitmaps",
        comment="保存中国 A 股各年的交易日真值位图。",
    ),
    "us_stock": BitmapSpec(
        calendar_id="us_stock",
        truth_table_path=ROOT / "tests" / "testdata" / "us_stock_truth_table.json",
        output_path=ROOT / "core" / "data" / "us_days_bitmaps.go",
        go_var_name="USStockTradeBitmaps",
        comment="保存美股各年的交易日真值位图。",
    ),
}


def load_truth_table(path: Path, expected_calendar_id: str) -> tuple[date, date, dict[str, bool]]:
    payload = json.loads(path.read_text(encoding="utf-8"))

    if payload.get("calendar_id") != expected_calendar_id:
        raise ValueError(f"unexpected calendar_id: {payload.get('calendar_id')}")

    start = date.fromisoformat(payload["start"])
    end = date.fromisoformat(payload["end"])
    days = payload["days"]
    if not isinstance(days, dict):
        raise ValueError("truth table days must be a JSON object")

    return start, end, days


def build_year_bitmap(year: int, days: dict[str, bool]) -> list[int]:
    bits = [0] * 6
    for day, is_trade_day in days.items():
        if not is_trade_day:
            continue

        current_day = date.fromisoformat(day)
        if current_day.year != year:
            continue

        bit_index = current_day.timetuple().tm_yday - 1
        word_index = bit_index // 64
        bit_offset = bit_index % 64
        if word_index >= len(bits):
            raise ValueError(f"day out of bitmap range: {day}")

        bits[word_index] |= 1 << bit_offset

    return bits


def render_go_source(spec: BitmapSpec, start: date, end: date, days: dict[str, bool]) -> str:
    years = list(range(start.year, end.year + 1))
    lines: list[str] = [
        "package data",
        "",
        f"// {spec.go_var_name} {spec.comment}",
        f"var {spec.go_var_name} = YearTradeBitmaps{{",
    ]

    for year in years:
        bits = build_year_bitmap(year, days)
        lines.append(f"\t{year}: {{")
        lines.append(f"\t\tYear: {year},")
        lines.append("\t\tBits: [6]uint64{")
        for value in bits:
            lines.append(f"\t\t\t{value},")
        lines.append("\t\t},")
        lines.append("\t},")

    lines.extend(
        [
            "}",
            "",
            "func init() {",
            f"\t{spec.go_var_name}.mustAlignYearKeys()",
            "}",
            "",
        ]
    )

    return "\n".join(lines)


def generate_bitmap_source(spec: BitmapSpec) -> str:
    start, end, days = load_truth_table(spec.truth_table_path, spec.calendar_id)
    return render_go_source(spec, start, end, days)


def write_bitmap_file(spec: BitmapSpec) -> Path:
    source = generate_bitmap_source(spec)
    with spec.output_path.open("w", encoding="utf-8", newline="\n") as file:
        file.write(source)
    return spec.output_path
