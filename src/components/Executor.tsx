import { spawn } from 'child_process';

/**
 * Result interface for shell command execution
 */
export interface CommandResult {
  output: string;
  error: string;
  returnCode: number;
}

/**
 * Executes a shell command and returns the result with output, error, and returnCode.
 *
 * @param command - The command to execute (e.g., 'ls', 'npm install')
 * @param args - Optional array of command arguments
 * @param options - Optional execution options
 * @returns Promise that resolves with CommandResult containing output, error, and returnCode
 *
 * @example
 * ```typescript
 * const result = await executeCommand('ls', ['-la']);
 * if (result.returnCode === 0) {
 *   console.log(result.output);
 * } else {
 *   console.error(result.error);
 * }
 * ```
 */
export async function executeCommand(
  command: string,
  args: string[] = [],
  options: {
    cwd?: string;
    env?: NodeJS.ProcessEnv;
    timeout?: number;
  } = {}
): Promise<CommandResult> {
  return new Promise((resolve) => {
    let output = '';
    let error = '';

    const { cwd, env, timeout } = options;

    const childProcess = spawn(command, args, {
      cwd,
      env: env || process.env,
      shell: process.platform === 'win32',
    });

    // Set timeout if provided
    if (timeout && timeout > 0) {
      setTimeout(() => {
        childProcess.kill('SIGTERM');
        error = `Command timed out after ${timeout}ms`;
        resolve({
          output,
          error,
          returnCode: -1,
        });
      }, timeout);
    }

    // Capture stdout
    childProcess.stdout?.on('data', (data: Buffer) => {
      output += data.toString();
    });

    // Capture stderr
    childProcess.stderr?.on('data', (data: Buffer) => {
      error += data.toString();
    });

    // Handle process completion
    childProcess.on('close', (code: number | null) => {
      resolve({
        output: output.trim(),
        error: error.trim(),
        returnCode: code !== null ? code : -1,
      });
    });

    // Handle process errors
    childProcess.on('error', (err: Error) => {
      resolve({
        output,
        error: err.message,
        returnCode: -1,
      });
    });
  });
}

/**
 * Convenience function to execute a command string (splits on spaces).
 * Useful for simple commands like "git status" or "npm test".
 *
 * @param commandString - The command string to execute
 * @param options - Optional execution options
 * @returns Promise that resolves with CommandResult
 *
 * @example
 * ```typescript
 * const result = await executeCommandString('git status');
 * ```
 */
export async function executeCommandString(
  commandString: string,
  options?: {
    cwd?: string;
    env?: NodeJS.ProcessEnv;
    timeout?: number;
  }
): Promise<CommandResult> {
  const [command, ...args] = commandString.split(' ');
  if (!command) {
    return {
      output: '',
      error: 'No command specified',
      returnCode: -1,
    };
  }
  return executeCommand(command, args.filter((arg): arg is string => arg !== undefined), options);
}
