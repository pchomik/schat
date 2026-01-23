import { executeCommand } from '../utils/ShellExecutor';

/**
 * Result interface for provider execution
 */
export interface ProviderResult {
    value: string;
    error?: string;
    success: boolean;
}

/**
 * CursorProvider executes cursor-agent shell commands with optional provider prompt prefix.
 */
export class CursorProvider {
    private providerPrompt: string = 'Return output in markdown format. Rationale is not needed.';

    /**
     * Executes cursor-agent with the given prompt.
     *
     * @param prompt - The user prompt to execute
     * @returns Promise resolving to ProviderResult with command output
     */
    async run(
        prompt: string,
        newSession: boolean = true,
        options?: {
            timeout?: number;
        }
    ): Promise<ProviderResult> {
        const fullPrompt = newSession ? `${this.providerPrompt} ${prompt}` : prompt;
        const commandArgs = [
            '--sandbox',
            'enabled',
            '--mode',
            'ask',
            ...(newSession ? [] : ['--continue']),
            '--model',
            'gemini-3-flash',
            '-p',
            `'${fullPrompt}'`,
        ];

        const result = await executeCommand('cursor-agent', commandArgs, {
            timeout: options?.timeout,
            cwd: process.cwd(),
        });

        if (result.returnCode === 0) {
            return {
                value: result.output,
                success: true,
            };
        }

        return {
            value: result.output,
            error: result.error || `Command failed with code ${result.returnCode}`,
            success: false,
        };
    }
}
