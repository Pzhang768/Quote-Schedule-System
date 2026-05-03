import { Quote } from "@/api/quotes";

interface Props {
  quote: Quote;
  selected: boolean;
  onClick: () => void;
}

export default function QuoteCard({ quote, selected, onClick }: Props) {
  return (
    <article
      className={`border rounded-xl p-4 cursor-pointer transition-colors ${selected ? "border-ink bg-ink/5" : "border-divider hover:bg-hover"}`}
      onClick={onClick}
    >
      <h3 className="text-heading">{quote.CustomerName}</h3>
      <p className="text-body text-muted">{quote.Address}</p>
      {quote.Description && <p className="text-caption text-muted mt-1">{quote.Description}</p>}
    </article>
  );
}
