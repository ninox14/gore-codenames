import type { GameState } from '@/types';
import { useEffect, useRef, useState, useCallback } from 'react';

export const teamColors = ['blue', 'red'] as const;
export type TeamColors = (typeof teamColors)[number];

// --- Server -> Client events ---
export interface ServerEventMap {
  // join_game: { playerId: string; name: string };
  // leave_game: { playerId: string };
  game_state: GameState;
  error: { message: string; err: string };
}

// --- Client -> Server events ---
export interface ClientEventMap {
  join_game: undefined;
  leave_game: { playerId: string };
  request_state: { gameId: string };
  change_team: {
    destination: string;
    // destination_color: TeamColors
  };
}

// --- Utility type: build discriminated union automatically ---
type EventUnion<TMap> = {
  [K in keyof TMap]: {
    type: K;
    game_id?: string;
    data: TMap[K];
  };
}[keyof TMap];

// --- Final types ---
export type ServerMessage = EventUnion<ServerEventMap>;
export type ClientMessage = EventUnion<ClientEventMap>;

type Options = {
  reconnect?: boolean;
  reconnectInterval?: number;
  onOpenCallback?: (socket: WebSocket) => void;
};

// Generic hook: TIn = client->server, TOut = server->client
export function useWebSocket<
  TIn extends { type: string },
  TOut extends { type: string }
>(url: string, options: Options = {}) {
  const {
    reconnect = true,
    reconnectInterval = 3000,

    onOpenCallback,
  } = options;

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeout = useRef<NodeJS.Timeout | null>(null);

  const [isConnected, setIsConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<TOut | null>(null);

  const connect = useCallback(() => {
    if (wsRef.current) return;

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      onOpenCallback?.(ws);
    };

    ws.onclose = () => {
      setIsConnected(false);
      wsRef.current = null;
      if (reconnect) {
        reconnectTimeout.current = setTimeout(connect, reconnectInterval);
      }
    };

    ws.onerror = (err) => {
      console.error('WebSocket error:', err);
      ws.close();
    };

    ws.onmessage = (event) => {
      try {
        const data: TOut = JSON.parse(event.data);
        console.log(JSON.parse(event.data));
        setLastMessage(data);
      } catch {
        console.warn('Invalid WS message:', event.data);
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [url, reconnect, reconnectInterval]);

  // send is now strictly typed with TIn
  const send = useCallback((msg: TIn) => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(msg));
    } else {
      console.warn('WebSocket not connected');
    }
  }, []);

  const close = useCallback(() => {
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.close();
    } else {
      console.warn('WebSocket not connected');
    }
  }, []);

  useEffect(() => {
    connect();
    return () => {
      if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
      wsRef.current?.close();
      wsRef.current = null;
    };
  }, [connect]);

  return { isConnected, lastMessage, send, close };
}
