import { useNavigate } from "react-router-dom";
import { FaCheckCircle, FaLock, FaClipboardList, FaBell } from "react-icons/fa";

export default function HomePage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Hero Section */}
      <section className="relative bg-blue-600 text-white text-center py-20">
        <div className="max-w-4xl mx-auto px-6">
          <h1 className="text-4xl font-bold mb-4">
            Revolutionizing Result Distribution at Tribhuvan University
          </h1>
          <p className="text-lg mb-6">
            Get instant access to your academic results, exam schedules, and more!
          </p>
          <div className="flex justify-center space-x-4">
            <button
              onClick={() => navigate("/login")}
              className="bg-white text-blue-600 px-6 py-3 rounded-lg font-semibold shadow-lg hover:bg-gray-200 transition"
            >
              Login
            </button>
            <button
              onClick={() => navigate("/register")}
              className="bg-blue-700 px-6 py-3 rounded-lg font-semibold shadow-lg hover:bg-blue-800 transition"
            >
              Sign Up
            </button>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-white">
        <div className="max-w-5xl mx-auto px-6">
          <h2 className="text-3xl font-bold text-center mb-10">Why Use Result-E?</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            <FeatureCard
              icon={<FaCheckCircle className="text-blue-600 text-3xl" />}
              title="Real-time Access to Results"
              description="View your academic results as soon as they are published."
            />
            <FeatureCard
              icon={<FaClipboardList className="text-green-600 text-3xl" />}
              title="Efficient Student Management"
              description="Administrators can easily manage student records, courses, and marks."
            />
            <FeatureCard
              icon={<FaLock className="text-red-600 text-3xl" />}
              title="Secure & Reliable"
              description="Encrypted login & secure database management for all users."
            />
            <FeatureCard
              icon={<FaBell className="text-yellow-600 text-3xl" />}
              title="Instant Notifications"
              description="Get notified about exam schedules, results, and announcements."
            />
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-16 bg-gray-100">
        <div className="max-w-4xl mx-auto px-6">
          <h2 className="text-3xl font-bold text-center mb-10">How It Works</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 text-center">
            <StepCard number="1" title="Login to Result-E" description="Secure access for students & admins." />
            <StepCard number="2" title="View Exam & Results" description="Instant access to published results." />
            <StepCard number="3" title="Stay Updated" description="Get notified about exams & notices." />
          </div>
        </div>
      </section>

      {/* Call to Action */}
      <section className="py-16 bg-blue-600 text-white text-center">
        <div className="max-w-3xl mx-auto px-6">
          <h2 className="text-3xl font-bold mb-6">Ready to Access Your Results?</h2>
          <p className="text-lg mb-6">
            Sign up now or log in to check your academic progress and stay updated!
          </p>
          <button
            onClick={() => navigate("/login")}
            className="bg-white text-blue-600 px-6 py-3 rounded-lg font-semibold shadow-lg hover:bg-gray-200 transition"
          >
            Get Started
          </button>
        </div>
      </section>
    </div>
  );
}

// Feature Card Component
const FeatureCard = ({ icon, title, description }) => (
  <div className="bg-gray-50 p-6 rounded-lg shadow-md text-center">
    <div className="mb-4">{icon}</div>
    <h3 className="text-lg font-semibold mb-2">{title}</h3>
    <p className="text-gray-600">{description}</p>
  </div>
);

// Step Card Component
const StepCard = ({ number, title, description }) => (
  <div className="bg-white p-6 rounded-lg shadow-md">
    <span className="text-3xl font-bold text-blue-600">{number}</span>
    <h3 className="text-lg font-semibold mt-3">{title}</h3>
    <p className="text-gray-600 mt-2">{description}</p>
  </div>
);
