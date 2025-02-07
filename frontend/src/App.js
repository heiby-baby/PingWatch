import React, { useEffect, useState } from 'react';
import './App.css';

function App() {
  const [results, setResults] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8080/ping-results')
      .then((response) => response.json())
      .catch((error) => console.error('Error fetching data:', error));
  }, []);

  const formatDate = (dateString) => {
    const date = new Date(dateString);
    return date.toLocaleString();   
  };

  return (
    <div className="App">
      <h1>Ping Results</h1>
      <table>
        <thead>
          <tr>
            <th>IP Address</th>
            <th>Ping Time</th>
            <th>Last Success</th>
          </tr>
        </thead>
        <tbody>
          {results.map((result, index) => (
            <tr key={index}>
              <td>{result.ip}</td>
              <td>{formatDate(result.ping_time)}</td>
              <td>{result.last_success ? formatDate(result.last_success) : 'N/A'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default App;