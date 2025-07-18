import React, { useState } from 'react';
import axios from '../axiosConfig';

const RecreateIndexButton = ({ endpoint, listName, onSuccess, onError }) => {
  const [isLoading, setIsLoading] = useState(false);
  const [message, setMessage] = useState('');

  const handleRecreateIndex = async () => {
    if (!window.confirm(`Are you sure you want to recreate the ${listName} index? This will delete and rebuild the Elasticsearch index with all current data.`)) {
      return;
    }

    setIsLoading(true);
    setMessage('');

    try {
      const response = await axios.post(endpoint);
      setMessage(response.data.message);
      if (onSuccess) {
        onSuccess(response.data.message);
      }
    } catch (error) {
      const errorMessage = error.response?.data?.error || 'Failed to recreate index';
      setMessage(`Error: ${errorMessage}`);
      if (onError) {
        onError(errorMessage);
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="recreate-index-container">
      <button
        onClick={handleRecreateIndex}
        disabled={isLoading}
        className={`btn btn-warning ${isLoading ? 'loading' : ''}`}
        title={`Recreate ${listName} Elasticsearch index`}
      >
        {isLoading ? (
          <>
            <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
            Recreating...
          </>
        ) : (
          <>
            <i className="fas fa-sync-alt me-2"></i>
            Recreate Index
          </>
        )}
      </button>
      {message && (
        <div className={`alert ${message.startsWith('Error') ? 'alert-danger' : 'alert-success'} mt-2`}>
          {message}
        </div>
      )}
    </div>
  );
};

export default RecreateIndexButton; 