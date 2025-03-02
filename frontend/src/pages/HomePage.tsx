import { Link } from 'react-router-dom';

export const HomePage = () => {
  return (
    <div className="container mx-auto">
      <header className="App-header">
        <h1 className="text-2xl font-bold text-center">React Frontend</h1>
        <p className="text-lg text-center">Go Backend</p>
      </header>
      <div className="text-center mt-4">
        <Link to="/login" className="text-blue-500 hover:text-blue-700">
          Go to Login
        </Link>
      </div>
    </div>
  );
}