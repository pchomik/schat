import React, { useMemo } from 'react';
import { Box, Text } from 'ink';
import { marked } from 'marked';
// @ts-expect-error - no types available
import { markedTerminal } from 'marked-terminal';

marked.use(markedTerminal());

interface ResponseProps {
    value: string
}

export const Response: React.FC<ResponseProps> = React.memo(({ value }: ResponseProps) => {
    const rendered = useMemo(() => {
        return marked.parse(value, { async: false }) as string;
    }, [value]);

    return (
        <Box padding={1} flexDirection="column">
            <Text>{rendered}</Text>
        </Box>
    );
});
