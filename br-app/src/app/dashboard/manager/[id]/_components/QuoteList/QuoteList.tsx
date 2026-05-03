import { Quote } from "@/api/quotes";
import Pagination from "@/components/Pagination/Pagination";
import QuoteCard from "../QuoteCard/QuoteCard";

interface Props {
  quotes: Quote[];
  selectedQuote: Quote | null;
  onSelect: (quote: Quote) => void;
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
}

export default function QuoteList({
  quotes,
  selectedQuote,
  onSelect,
  page,
  totalPages,
  onPrev,
  onNext,
}: Props) {
  return (
    <section className="flex-2 flex flex-col gap-2">
      <h2 className="text-caption-upper text-muted px-1">Unscheduled Quotes</h2>
      {quotes.length === 0 && (
        <div className="text-body text-muted p-2">No unscheduled quotes.</div>
      )}
      {quotes.map((q) => (
        <QuoteCard
          key={q.ID}
          quote={q}
          selected={selectedQuote?.ID === q.ID}
          onClick={() => onSelect(q)}
        />
      ))}
      <Pagination page={page} totalPages={totalPages} onPrev={onPrev} onNext={onNext} />
    </section>
  );
}
