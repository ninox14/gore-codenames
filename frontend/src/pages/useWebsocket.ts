import { useEffect, useRef, useState } from 'react';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function handleMessage(msg: any) {
  console.log(msg);
}

function useGameWS(wsUrl: string) {
  const wsRef = useRef<WebSocket | null>(null);
  const [status, setStatus] = useState<string>('');

  useEffect(() => {
    const ws = new WebSocket(wsUrl);
    ws.onopen = function () {
      setStatus('Connected to server');
    };

    ws.onmessage = function (event) {
      console.log(event);
      const message = JSON.parse(event.data);
      handleMessage(message);
    };

    ws.onclose = function () {
      setStatus('Disconnected from server');
      wsRef.current = null;
    };

    ws.onerror = function (error) {
      setStatus('Connection error');
      console.error('WebSocket error:', error);
    };

    wsRef.current = ws;
  }, [wsUrl]);

  return { ws: wsRef.current, status };
}

export default useGameWS;
