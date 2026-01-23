import React from 'react';
import { Box, Text, Static } from 'ink';
import { useInput } from 'ink';
import { Prompt } from './ui/Prompt';
import { TextArea } from './ui/TextArea';
import { Response } from './ui/Response';
import { CursorProvider } from './providers/Cursor';
import Spinner from 'ink-spinner';

interface Exchange {
    id: number;
    prompt: string;
    response: string;
}

export const App: React.FC = () => {
    const [history, setHistory] = React.useState<Exchange[]>([]);
    const [currentPrompt, setCurrentPrompt] = React.useState<string | null>(null);
    const [isProcessing, setIsProcessing] = React.useState(false);
    const cursorProviderRef = React.useRef(new CursorProvider());
    const [newSession, setNewSession] = React.useState(true);
    const idRef = React.useRef(0);

    const addNewPrompt = React.useCallback((newPrompt: string) => {
        const trimmedPrompt = newPrompt.trimEnd();
        if (!trimmedPrompt.trim()) {
            return;
        }

        const exchangeId = idRef.current++;
        setCurrentPrompt(trimmedPrompt);
        setIsProcessing(true);

        cursorProviderRef.current.run(trimmedPrompt, newSession).then((result: any) => {
            const responseText = result?.success
                ? (result.value || '(no output)')
                : `Error: ${result?.error || result?.value || 'Unknown error'}`;

            setHistory((prev) => [...prev, {
                id: exchangeId,
                prompt: trimmedPrompt,
                response: responseText,
            }]);
            setCurrentPrompt(null);
            setNewSession(false);
            setIsProcessing(false);
        }).catch((err: any) => {
            setHistory((prev) => [...prev, {
                id: exchangeId,
                prompt: trimmedPrompt,
                response: `Error: ${err?.message || err}`,
            }]);
            setCurrentPrompt(null);
            setNewSession(false);
            setIsProcessing(false);
        });
    }, [newSession]);

    return (
        <Box flexDirection="column">
            <Static items={history}>
                {(exchange) => (
                    <Box key={exchange.id} flexDirection="column">
                        <Prompt value={exchange.prompt} />
                        <Response value={exchange.response} />
                    </Box>
                )}
            </Static>
            {currentPrompt && <Prompt value={currentPrompt} />}
            {isProcessing && <Box flexDirection="row"><Spinner type="dots" /><Text color="rgb(131,139,167)"> Processing...</Text></Box>}
            <TextArea newPromptCallback={addNewPrompt} />
            <Text color="rgb(131,139,167)">| CTRL + S - send | CTRL + C - exit |</Text>
        </Box>
    )
};
