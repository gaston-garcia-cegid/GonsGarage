// src/components/client/ClientCars.tsx
export default function ClientCars(
  { onAddCar }: { onAddCar: () => void }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-lg font-semibold">Your Vehicles</h3>
        <button onClick={onAddCar}
          className="bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700">
          Add New Car
        </button>
      </div>
      <p>Your cars will be listed here...</p>
    </div>
  );
}