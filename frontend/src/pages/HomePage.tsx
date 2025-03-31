import { Link } from 'react-router-dom';
import { Navbar } from '../components/layout/Navbar';

export const HomePage = () => {
  return (
    <div className="w-full">
      <Navbar/>
        <div className="hero bg-blue-100 flex flex-col justify-center items-center h-screen">
          <h1 className="text-4xl font-bold text-center text-blue-500">Welcome to Our Site</h1>
          <p className="text-lg text-center mt-4">Your journey to knowledge starts here.</p>
        </div>
      <div className="text-center mt-4">
        <Link to="/login" className="text-blue-500 hover:text-blue-700">
          Go to Login
        </Link>
      </div>
    </div>
  );
}