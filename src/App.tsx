import React from 'react';
import { Box, Text } from 'ink';
import { Prompt } from './components/Prompt';
import { TextArea } from './components/TextArea';

export const App: React.FC = () => {
    const [prompts, setPrompts] = React.useState<string[]>([]);

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
            <TextArea newPromptCallback={addNewPrompt} resetAppCallback={resetApp}/>
            <Text color="rgb(131,139,167)">| CTRL + S - send | CTRL + N - new session | CTRL + C - exit |</Text>
        </Box>
    )
};
