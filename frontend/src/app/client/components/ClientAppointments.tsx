// src/components/client/ClientAppointments.tsx
export default function ClientAppointments(
  { onScheduleService }: { onScheduleService: () => void }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-lg font-semibold">Your Appointments</h3>
        <button 
          onClick={onScheduleService}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
          Schedule Service
        </button>
      </div>
      <p>Your appointments will be listed here...</p>
    </div>
  );
}