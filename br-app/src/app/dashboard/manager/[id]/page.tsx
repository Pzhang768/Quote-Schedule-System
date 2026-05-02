"use client";
import { useCallback, useState } from "react";
import { useParams } from "next/navigation";
import usePaginated from "@/hooks/usePaginated/usePaginated";
import { getQuotes, Quote } from "@/api/quotes";
import { getTechnicians, Technician } from "@/api/technicians";
import QuoteList from "./_components/QuoteList/QuoteList";
import TechnicianList from "./_components/TechnicianList/TechnicianList";
import AssignForm from "./_components/AssignForm/AssignForm";

const today = new Date().toLocaleDateString("en-CA");

interface Selection {
  quote: Quote | null;
  technician: Technician | null;
  time: string;
}

const emptySelection: Selection = { quote: null, technician: null, time: "" };

export default function ManagerDashboard() {
  const { id: managerId } = useParams<{ id: string }>();

  const quotesFetcher = useCallback(
    (page: number, ps: number) => getQuotes(page, ps, "unscheduled"),
    []
  );
  const techFetcher = useCallback((page: number, ps: number) => getTechnicians(page, ps), []);

  const {
    data: quotes,
    refetch,
    page: quotePage,
    totalPages: quoteTotalPages,
    nextPage: quotesNext,
    prevPage: quotesPrev,
  } = usePaginated<Quote>(quotesFetcher);
  const {
    data: technicians,
    page: techPage,
    totalPages: techTotalPages,
    nextPage: techNext,
    prevPage: techPrev,
  } = usePaginated<Technician>(techFetcher);

  const [selection, setSelection] = useState<Selection>(emptySelection);
  const [date, setDate] = useState(today);

  const canAssign = !!selection.quote && !!selection.technician && !!selection.time;

  const handleSuccess = () => {
    setSelection(emptySelection);
    refetch();
  };

  return (
    <div className="flex flex-col h-full w-full p-4 gap-4">
      <div className="border border-divider rounded-xl px-6 py-4 flex items-center">
        <h1 className="text-title">Manager Dashboard</h1>
      </div>

      <div className="flex gap-4 flex-1">
        <QuoteList
          quotes={quotes}
          selectedQuote={selection.quote}
          onSelect={(q) => setSelection({ ...emptySelection, quote: q })}
          page={quotePage}
          totalPages={quoteTotalPages}
          onPrev={quotesPrev}
          onNext={quotesNext}
        />
        <TechnicianList
          technicians={technicians}
          quoteSelected={!!selection.quote}
          date={date}
          onDateChange={setDate}
          selectedTechnicianId={selection.technician?.ID ?? null}
          selectedTime={selection.time}
          onSelect={(t) => setSelection((s) => ({ ...s, technician: t, time: "" }))}
          onTimeSelect={(t, time) => setSelection((s) => ({ ...s, technician: t, time }))}
          page={techPage}
          totalPages={techTotalPages}
          onPrev={techPrev}
          onNext={techNext}
        />
      </div>

      {canAssign && selection.quote && selection.technician && (
        <AssignForm
          managerId={managerId}
          quote={selection.quote}
          technician={selection.technician}
          date={date}
          time={selection.time}
          onSuccess={handleSuccess}
        />
      )}
    </div>
  );
}
