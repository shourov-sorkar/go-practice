import { Link } from "react-router-dom";

export const Navbar = () => {
  return (
    <nav className="bg-white shadow-lg">
      <div className="max-w-6xl mx-auto px-4">
        <div className="flex justify-between">
          <div className="flex space-x-7">
            <div className="flex items-center py-4 px-2">
              <h1 className="text-gray-800 text-2xl font-bold">
                <Link to="/">Home</Link>
              </h1>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};
