import React from 'react';
import { Box, Text } from 'ink';
import { useInput } from 'ink';
import { Prompt } from './components/Prompt';
import { TextArea } from './components/TextArea';

export const App: React.FC = () => {
    const [prompts, setPrompts] = React.useState<string[]>([]);

    useInput((input, key) => {
        if (key.ctrl && input === 'n') {
            setPrompts([]);
        }
    });

    const addNewPrompt = (newPrompt: string) => {
        setPrompts((prev) => [...prev, newPrompt]);
    }

    const resetApp = () => {
        setPrompts([]);
    }

	return (
        <Box flexDirection="column">
            {prompts.map((prompt, index) => (
                <Prompt key={index} value={prompt} />
            ))}
            <TextArea newPromptCallback={addNewPrompt}/>
            <Text color="rgb(131,139,167)">| CTRL + S - send | CTRL + N - new session | CTRL + C - exit |</Text>
        </Box>
    )
};
