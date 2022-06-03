#!/usr/bin/env python3

import logging
import sys
import argparse
import distutils.spawn
import subprocess
import os

logger = logging.getLogger()

def setupLogging(log_verbosity):
    logger.setLevel(getattr(logging, log_verbosity))
    handler = logging.StreamHandler(sys.stdout)
    formatter = logging.Formatter(
        fmt = '[%(levelname)s] (%(asctime)s): %(message)s',
        datefmt = '%m/%d/%Y %I:%M:%S %p')
    handler.setFormatter(formatter)
    logger.addHandler(handler)


def tagAndPushImage(dockerpath, image_name, registry, repo, tag):
    cmd = "%s images -q %s" % (dockerpath, image_name)
    logging.debug("Executing: [%s]" % cmd)
    pp = subprocess.Popen(
        cmd, 
        shell=True, 
        stdout=subprocess.PIPE)
    image_id = pp.stdout.read().strip().decode("utf-8")
    logging.debug("we got image id  [%s]" % image_id)
    
    cmd = "%s tag %s %s/%s:%s" % (dockerpath, image_id, registry, repo, tag)
    logging.debug("Executing: [%s]" % cmd)
    pp = subprocess.Popen(
        cmd, 
        shell=True, 
        stdout=subprocess.PIPE)
   
    out1 = pp.stdout.read().strip().decode("utf-8")
    logging.debug("out 1 : [%s]" % out1)

    cmd = "%s push %s/%s:%s" % (dockerpath, registry, repo, tag)
    logging.debug("Executing: [%s]" % cmd)
    pp = subprocess.Popen(
        cmd, 
        shell=True, 
        stdout=subprocess.PIPE)
    out2 = pp.stdout.read().strip().decode("utf-8")
    logging.debug("out 2 : [%s]" % out2)


def authToDocker(
        awsclipath,
        registry,
        aws_profile):

    try:
        p1 = subprocess.Popen(
            [awsclipath, "ecr", "get-login-password", ("--profile=%s" % aws_profile)],
            stdout = subprocess.PIPE)
    except:
        return False

    try:
        p2 = subprocess.Popen(
            ["docker",
            "login","--username", "AWS", "--password-stdin", registry],
             stdin = p1.stdout)

        p1.stdout.close()
        out, err = p2.communicate()
    except:
        return False   
             
    return True

def _main():
    parser = argparse.ArgumentParser(allow_abbrev=False)
    parser.add_argument('--registry', default='%{registry}')
    parser.add_argument('--repo', default='%{repo}')
    parser.add_argument('--image_name', default='%{image_name}')
    parser.add_argument('--aws_profile', default='%{aws_profile}')
    parser.add_argument('--log_verbosity', default='%{log_verbosity}')
    parser.add_argument('--extra_tag', default='%{extra_tag}')

    flags, args = parser.parse_known_args()
    setupLogging(flags.log_verbosity)

    logging.debug("We got registry: [%s]" % flags.registry)
    logging.debug("We got repo: [%s]" % flags.repo)
    logging.debug("We got image name: [%s]" % flags.image_name)
    logging.debug("We got aws profile: [%s]" % flags.aws_profile)
    logging.debug("We got extra tag: [%s]" % flags.extra_tag)

    extra_tag = flags.extra_tag.strip()
    if " " in extra_tag:
        print("Bad extra tag: [%s]" % extra_tag)
        sys.exit(-1)

    # TODO: Following section needs to auth against ticket
    awsclipath = distutils.spawn.find_executable("aws")
    if awsclipath == None:
        print("No aws client found, you need to install it")
        sys.exit(-1)
    
    dockerpath = distutils.spawn.find_executable("docker")

    if authToDocker(
        awsclipath,
        flags.registry,
        flags.aws_profile) == False:
        sys.exit(-2)

    gitpath = distutils.spawn.find_executable("git")
    build_working_dir = os.environ.get("BUILD_WORKING_DIRECTORY")
    logging.debug("CWD: [%s]" % build_working_dir)

    cmd = "cd %s && %s describe --tags $(%s rev-list --tags --max-count=1)" % (build_working_dir, gitpath, gitpath)    
    logging.debug("Executing to get latest tag: [%s]" % cmd)
    pp = subprocess.Popen(
        cmd, 
        shell=True, 
        stdout=subprocess.PIPE)
    latest_tag = pp.stdout.read().strip().decode("utf-8")

    logging.debug("Latest tag is %s" % latest_tag)

    cmd = "cd %s && %s tag --points-at HEAD" % (build_working_dir, gitpath)
    logging.debug("Executing to get current tag: [%s]" % cmd)
    pp = subprocess.Popen(
        cmd, 
        shell=True, 
        stdout=subprocess.PIPE)
    current_tag = pp.stdout.read().strip().decode("utf-8")
    logging.debug("Current tag is %s" % current_tag)

    tagAndPushImage(
        dockerpath, 
        flags.image_name, 
        flags.registry, 
        flags.repo,
        latest_tag)
    
    tagAndPushImage(
        dockerpath, 
        flags.image_name, 
        flags.registry, 
        flags.repo,
        extra_tag)
    
_main()
