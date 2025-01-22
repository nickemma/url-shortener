import React, { useState } from "react";
import axios from "axios";

interface Stats {
  original_url: string;
  short_code: string;
  click_count: number;
  created_at: string;
}

const App: React.FC = () => {
  const [originalURL, setOriginalURL] = useState<string>("");
  const [shortURL, setShortURL] = useState<string>("");
  const [stats, setStats] = useState<Stats | null>(null);

  const handleShorten = async (): Promise<void> => {
    try {
      const response = await axios.post<{ short_url: string }>(
        "http://localhost:5000/shorten", 
      {
        original_url: originalURL,
      },
     {
    headers: {
      "Content-Type": "application/json",
    },
  }
    );
      console.log("res", response)
      setShortURL(response.data.short_url);
      setStats(null); // Clear previous stats
    } catch (error) {
      console.error("Error creating short URL:", error);
    }
  };

  const fetchStats = async (): Promise<void> => {
    try {
      const shortCode = shortURL.split("/").pop() || ""; // Extract the short code from the URL
      const response = await axios.get<Stats>(`http://localhost:5000/stats/${shortCode}`);
      setStats(response.data);
    } catch (error) {
      console.error("Error fetching stats:", error);
    }
  };

  return (
    <div className="p-6 flex flex-col items-center">
      <h1 className="text-2xl font-bold mb-4">URL Shortener</h1>
      <div className="flex mb-4">
        <input
          type="text"
          placeholder="Enter URL"
          value={originalURL}
          onChange={(e) => setOriginalURL(e.target.value)}
          className="border p-2 w-64 rounded-l"
        />
        <button onClick={handleShorten} className="bg-blue-500 text-white p-2 rounded-r">
          Shorten
        </button>
      </div>
      {shortURL && (
        <div className="mb-4">
          <p>Shortened URL:</p>
          <a href={shortURL} target="_blank" rel="noopener noreferrer" className="text-blue-600">
            {shortURL}
          </a>
          <button onClick={fetchStats} className="ml-4 bg-gray-300 p-2 rounded">
            View Stats
          </button>
        </div>
      )}
      {stats && (
        <div className="mt-4 p-4 border rounded bg-gray-100">
          <h2 className="font-bold">URL Statistics</h2>
          <p>
            <strong>Original URL:</strong> {stats.original_url}
          </p>
          <p>
            <strong>Short Code:</strong> {stats.short_code}
          </p>
          <p>
            <strong>Clicks:</strong> {stats.click_count}
          </p>
          <p>
            <strong>Created At:</strong> {new Date(stats.created_at).toLocaleString()}
          </p>
        </div>
      )}
    </div>
  );
};

export default App;

