import { Button } from '@/components/ui/button';
import { Link } from 'react-router';

function Home() {
  return (
    <div className="min-h-screen bg-black py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-gray-50 mb-2">
            Welcome to GoRe Codenames
          </h1>
          <p className="text-gray-600">Get started by editing </p>

          <Button asChild>
            <Link className="mt-3" to="/game">
              Play
            </Link>
          </Button>
        </div>

        <div className="text-center text-gray-500 text-sm">
          Built with Vite, React, and Tailwind CSS
        </div>
      </div>
    </div>
  );
}

export default Home;
