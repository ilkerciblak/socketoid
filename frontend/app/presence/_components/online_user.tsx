export default function OnlineUser(user: { username: string; id: string }) {
  return (
    <li className="bg-black border border-gray-200 shadow-white shadow-sm flex flex-row px-2 py-1 rounded-md justify-around items-center w-55">
      <div className="rounded-full bg-green-700 shadow-green-200 shadow-md w-3 h-3"></div>
      
      <p className=" italic text-gray-100"> <span className="text-gray-400 mr-3">User name:</span>{user.username}</p>

    </li>
  );
}
