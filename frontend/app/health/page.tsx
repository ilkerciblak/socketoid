import { server } from "../lib/fetcher";

export default function Health() {
  type Health = {
    status: number;
    time: string;
  };
  const HealthStatus = async () => {
    const healthData = (await server.GET<Health>("health")) as Health;

    return (
      <div className="bg-blue-500 rounded-2xl flex flex-row justify-between px-1 py-2">
        <div className="bg-white rounded-full w-5 text-blue-500 italic">
          i
        </div>
        <p className="text-amber-400 italic"> Status: {healthData.status}</p>
        <p className="ml-2 text-sm text-amber-300"> at {healthData.time}</p>
      </div>
    );
  };

  return (
    <main>
      <HealthStatus />
    </main>
  );
}
