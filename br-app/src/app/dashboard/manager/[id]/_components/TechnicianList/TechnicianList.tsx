import { Technician } from "@/api/technicians";
import Pagination from "@/components/Pagination/Pagination";
import TechnicianRow from "../TechnicianRow/TechnicianRow";

interface Props {
  technicians: Technician[];
  quoteSelected: boolean;
  date: string;
  onDateChange: (date: string) => void;
  selectedTechnicianId: string | null;
  selectedTime: string;
  onSelect: (technician: Technician) => void;
  onTimeSelect: (technician: Technician, time: string) => void;
  page: number;
  totalPages: number;
  onPrev: () => void;
  onNext: () => void;
}

const today = () => new Date().toLocaleDateString("en-CA");

export default function TechnicianList({
  technicians,
  quoteSelected,
  date,
  onDateChange,
  selectedTechnicianId,
  selectedTime,
  onSelect,
  onTimeSelect,
  page,
  totalPages,
  onPrev,
  onNext,
}: Props) {
  return (
    <section className="flex-5 flex flex-col gap-2">
      <div className="flex items-center gap-3 px-1">
        <h2 className="text-caption-upper text-muted">Technicians</h2>
        <input
          type="date"
          value={date}
          min={today()}
          onChange={(e) => onDateChange(e.target.value)}
          className="text-caption border border-divider rounded-lg px-2 py-1 bg-transparent"
        />
      </div>
      {!quoteSelected && <div className="text-body text-muted p-2">Select a quote first.</div>}
      {quoteSelected &&
        technicians.map((t) => (
          <TechnicianRow
            key={t.ID}
            technician={t}
            date={date}
            selected={selectedTechnicianId === t.ID}
            onClick={() => onSelect(t)}
            selectedTime={selectedTechnicianId === t.ID ? selectedTime : ""}
            onTimeSelect={(time) => onTimeSelect(t, time)}
          />
        ))}
      <Pagination page={page} totalPages={totalPages} onPrev={onPrev} onNext={onNext} />
    </section>
  );
}
