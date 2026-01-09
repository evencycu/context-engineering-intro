import React, { useState, useEffect, useRef } from 'react';
import { service } from '../services';
import { UserProfile } from '../types';
import { Search, User, X } from 'lucide-react';

interface PeoplePickerProps {
  onSelect: (value: string) => void;
  value: string;
}

export const PeoplePicker: React.FC<PeoplePickerProps> = ({ onSelect, value }) => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<UserProfile[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showResults, setShowResults] = useState(false);
  const wrapperRef = useRef<HTMLDivElement>(null);

  // Debounce search
  useEffect(() => {
    const timer = setTimeout(async () => {
      if (query.length >= 2) {
        setIsLoading(true);
        try {
          const users = await service.searchUsers(query);
          setResults(users);
          setShowResults(true);
        } catch (error) {
          console.error("Search failed", error);
        } finally {
          setIsLoading(false);
        }
      } else {
        setResults([]);
        setShowResults(false);
      }
    }, 400);

    return () => clearTimeout(timer);
  }, [query]);

  // Click outside to close
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target as Node)) {
        setShowResults(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [wrapperRef]);

  const handleSelect = (user: UserProfile) => {
    onSelect(user.userPrincipalName);
    setQuery('');
    setShowResults(false);
  };

  const handleClear = () => {
    onSelect('');
    setQuery('');
  };

  return (
    <div className="relative w-full" ref={wrapperRef}>
      <label className="block text-sm font-medium text-gray-700 mb-1">Recipient (Azure AD Search)</label>
      
      {value ? (
        <div className="flex items-center justify-between p-2 border border-blue-200 bg-blue-50 rounded-md">
          <div className="flex items-center gap-2">
            <div className="bg-blue-100 p-1 rounded-full">
              <User size={16} className="text-blue-600" />
            </div>
            <span className="text-sm font-medium text-blue-900">{value}</span>
          </div>
          <button onClick={handleClear} className="text-blue-400 hover:text-blue-600">
            <X size={16} />
          </button>
        </div>
      ) : (
        <div className="relative">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search size={16} className="text-gray-400" />
          </div>
          <input
            type="text"
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            placeholder="Search by name or email..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
          />
          {isLoading && (
            <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
              <div className="animate-spin h-4 w-4 border-2 border-blue-500 rounded-full border-t-transparent"></div>
            </div>
          )}
        </div>
      )}

      {showResults && results.length > 0 && (
        <ul className="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
          {results.map((user) => (
            <li
              key={user.id}
              className="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-blue-50"
              onClick={() => handleSelect(user)}
            >
              <div className="flex flex-col">
                <span className="font-medium text-gray-900">{user.displayName}</span>
                <span className="text-gray-500 text-xs">{user.jobTitle} &bull; {user.userPrincipalName}</span>
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};