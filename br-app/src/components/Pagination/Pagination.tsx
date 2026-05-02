interface Props {
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
}

export default function Pagination({ page, totalPages, onPrev, onNext }: Props) {
  if (totalPages <= 1) return null;
  return (
    <div className="flex items-center gap-3 px-1 pt-1">
      <button
        onClick={onPrev}
        disabled={page <= 1}
        className="text-caption px-3 py-1 rounded-lg border border-divider disabled:opacity-40 hover:bg-hover transition-colors"
      >
        Prev
      </button>
      <span className="text-caption text-muted">
        {page} / {totalPages}
      </span>
      <button
        onClick={onNext}
        disabled={page >= totalPages}
        className="text-caption px-3 py-1 rounded-lg border border-divider disabled:opacity-40 hover:bg-hover transition-colors"
      >
        Next
      </button>
    </div>
  );
}
