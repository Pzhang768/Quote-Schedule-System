import { Quote } from "@/api/quotes";

interface Props {
  quote: Quote;
  selected: boolean;
  onClick: () => void;
}

export default function QuoteCard({ quote, selected, onClick }: Props) {
  return (
    <div
      className={`border rounded-xl p-4 cursor-pointer transition-colors ${selected ? "border-ink bg-ink/5" : "border-divider hover:bg-hover"}`}
      onClick={onClick}
    >
      <div className="text-heading">{quote.CustomerName}</div>
      <div className="text-body text-muted">{quote.Address}</div>
      {quote.Description && <div className="text-caption text-muted mt-1">{quote.Description}</div>}
    </div>
  );
}
