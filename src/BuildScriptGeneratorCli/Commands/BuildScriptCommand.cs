﻿// --------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.
// --------------------------------------------------------------------------------------------

using System;
using System.IO;
using System.Linq;
using McMaster.Extensions.CommandLineUtils;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Options;
using Microsoft.Oryx.BuildScriptGenerator;
using Microsoft.Oryx.BuildScriptGenerator.Common;
using Microsoft.Oryx.BuildScriptGenerator.Exceptions;
using Microsoft.Oryx.Common.Extensions;

namespace Microsoft.Oryx.BuildScriptGeneratorCli
{
    [Command(Name, Description = "Generate build script to standard output.")]
    internal class BuildScriptCommand : BuildCommandBase
    {
        public const string Name = "build-script";

        [Option(
            "--output",
            CommandOptionType.SingleValue,
            Description = "The path that the build script will be written to. " +
                          "If not specified, the result will be written to STDOUT.")]
        public string OutputPath { get; set; }

        internal override int Execute(IServiceProvider serviceProvider, IConsole console)
        {
            var scriptGenerator = new BuildScriptGenerator(
                serviceProvider,
                console,
                checkerMessageSink: null,
                operationId: null);

            if (!scriptGenerator.TryGenerateScript(out var generatedScript, out var exception))
            {
                if (exception != null)
                {
                    return ProcessExitCodeHelper.GetExitCodeForException(exception);
                }

                return ProcessConstants.ExitFailure;
            }

            if (string.IsNullOrEmpty(this.OutputPath))
            {
                console.WriteLine(generatedScript);
            }
            else
            {
                this.OutputPath.SafeWriteAllText(generatedScript);
                console.WriteLine($"Script written to '{this.OutputPath}'");

                // Try making the script executable
                ProcessHelper.TrySetExecutableMode(this.OutputPath);
            }

            return ProcessConstants.ExitSuccess;
        }

        internal override bool IsValidInput(IServiceProvider serviceProvider, IConsole console)
        {
            var options = serviceProvider.GetRequiredService<IOptions<BuildScriptGeneratorOptions>>().Value;

            if (!Directory.Exists(options.SourceDir))
            {
                console.WriteErrorLine($"Could not find the source directory '{options.SourceDir}'.");
                return false;
            }

            // Invalid to specify platform version without platform name
            if (string.IsNullOrEmpty(options.PlatformName) && !string.IsNullOrEmpty(options.PlatformVersion))
            {
                console.WriteErrorLine("Cannot use platform version without platform name also.");
                return false;
            }

            return true;
        }

        internal override void ConfigureBuildScriptGeneratorOptions(BuildScriptGeneratorOptions options)
        {
            BuildScriptGeneratorOptionsHelper.ConfigureBuildScriptGeneratorOptions(
                options,
                sourceDir: this.SourceDir,
                destinationDir: null,
                intermediateDir: null,
                manifestDir: null,
                platform: this.PlatformName,
                platformVersion: this.PlatformVersion,
                shouldPackage: this.ShouldPackage,
                requiredOsPackages: string.IsNullOrWhiteSpace(this.OsRequirements) ? null : this.OsRequirements.Split(','),
                appType: this.AppType,
                scriptOnly: true,
                properties: this.Properties);
        }

        internal override IServiceProvider TryGetServiceProvider(IConsole console)
        {
            // Don't use the IConsole instance in this method -- override this method in the command
            // and pass IConsole through to ServiceProviderBuilder to write to the output.
            var serviceProviderBuilder = new ServiceProviderBuilder(this.LogFilePath)
                .ConfigureServices(services =>
                {
                    var configuration = new ConfigurationBuilder()
                        .AddEnvironmentVariables()
                        .Build();

                    services.AddSingleton<IConfiguration>(configuration);
                })
                .ConfigureScriptGenerationOptions(opts =>
                {
                    this.ConfigureBuildScriptGeneratorOptions(opts);

                    // For debian flavor, we first check for existance of an environment variable
                    // which contains the os type. If this does not exist, parse the
                    // FilePaths.OsTypeFileName file for the correct flavor
                    if (string.IsNullOrWhiteSpace(opts.DebianFlavor))
                    {
                        var ostypeFilePath = Path.Join("/opt", "oryx", FilePaths.OsTypeFileName);
                        if (File.Exists(ostypeFilePath))
                        {
                            // these file contents are in the format <OS_type>|<Os_version>, e.g. DEBIAN|BULLSEYE
                            // we want the Os_version part only, as all lowercase
                            var fullOsTypeFileContents = File.ReadAllText(ostypeFilePath);
                            opts.DebianFlavor = fullOsTypeFileContents.Split("|").TakeLast(1).SingleOrDefault().Trim().ToLowerInvariant();
                        }
                        else
                        {
                            // If we cannot resolve the debian flavor, error out as we will not be able to determine
                            // the correct SDKs to pull
                            var errorMessage = $"Error: Image debian flavor not found in DEBIAN_FLAVOR environment variable or the " +
                                $"{Path.Join("/opt", "oryx", FilePaths.OsTypeFileName)} file. Exiting...";
                            throw new InvalidUsageException(errorMessage);
                        }
                    }
                });
            return serviceProviderBuilder.Build();
        }
    }
}