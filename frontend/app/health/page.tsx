"use server";
import SSEStatusIndicator from "../components/sse_status_component";
import { server } from "../lib/fetcher";

export default async function Health() {
  type Health = {
    status: number;
    time: string;
  };
  const HealthStatus = async () => {
    const healthData = (await server.GET<Health>("health")) as Health;

    return (
      <div className="w-100 bg-blue-500 rounded-2xl flex flex-row justify-between items-center px-3 m-auto py-2">
        <div className="bg-white rounded-xl align-bottom py-3 px-1 text-center text-blue-500 font-bold text-xs italic">Server Health</div>
        <p className="text-white font-bold italic"> Status: {healthData.status}</p>
        <p className="ml-2 text-sm text-amber-100"> at {healthData.time}</p>
      </div>
    );
  };

  return (
    <main>
      <HealthStatus />
      <SSEStatusIndicator/>
    </main>
  );
}
