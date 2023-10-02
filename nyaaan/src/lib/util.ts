export function getColor(asciiCode: number): string {
	let style = " font-semibold";
	switch (true) {
		case asciiCode >= 65 && asciiCode <= 68:
			return "text-success" + style;
		case asciiCode >= 69 && asciiCode <= 72:
			return "text-error" + style;
		case asciiCode >= 73 && asciiCode <= 77:
			return "text-accent" + style;
		case asciiCode >= 78 && asciiCode <= 81:
			return "text-secondary" + style;
		case asciiCode >= 82 && asciiCode <= 85:
			return "text-warning" + style;
		case asciiCode >= 86 && asciiCode <= 90:
			return "text-error" + style;
		default:
			return "text-accent" + style;
	}
}
